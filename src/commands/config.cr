module Soar::Commands
  class Config < Base
    def setup : Nil
      @name = "config"

      add_command Init.new
      add_option 'g', "global"
    end

    def run(arguments : Cling::Arguments, options : Cling::Options) : Nil
      config = Soar::Config.load
      puts config.to_yaml[4..]

      global = options.has? "global"
      if global
        config = Soar::Config.load_global
      else
        config = Soar::Config.load_local
      end

      if config.nil?
        stderr.puts "failed to load #{global ? "global" : "local"} config"
      else
        stdout.puts config.to_yaml[4..]
      end
    end

    private class Init < Base
      def setup : Nil
        @name = "init"

        add_argument "dir"
        add_option 'f', "force"
      end

      def run(arguments : Cling::Arguments, options : Cling::Options) : Nil
        path : Path

        if dir = arguments.get?("dir").try &.as_s
          unless Dir.exists? dir
            stderr.puts "directory does not exist"
            return
          end

          path = Path[dir, ".soar.yml"]
        else
          path = path = Soar::Config::PATH
          Dir.mkdir_p path.dirname
        end

        if File.file?(path) && !options.has?("force")
          stderr.puts "a config file already exists at this location"
          stderr.puts "re-run with the '--force' flag to overwrite"
          return
        end

        begin
          File.write path, Soar::Config.new.to_yaml[4..]
        rescue ex
          stderr.puts "failed to initialize config:"
          stderr.puts ex
        end
      end
    end
  end
end
