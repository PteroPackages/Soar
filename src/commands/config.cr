module Soar::Commands
  class Config < Base
    def setup : Nil
      @name = "config"

      add_command Init.new
      add_command Copy.new

      add_option 'g', "global"
      add_option 'l', "local"
    end

    def pre_run(arguments : Cling::Arguments, options : Cling::Options) : Bool
      if options.has?("global") && options.has?("local")
        error "cannot specify global and local option; pick one"
        return false
      end

      true
    end

    def run(arguments : Cling::Arguments, options : Cling::Options) : Nil
      global = options.has? "global"
      local = options.has? "local"

      config = if global
                 Soar::Config.load_global
               elsif local
                 Soar::Config.load_local
               else
                 Soar::Config.load
               end

      return error "failed to load #{global ? "global" : "local"} config" if config.nil?

      {% for field in %i(app client) %}
        stdout << "{{ field.id }} url: ".colorize.bold
        if config.{{ field.id }}?
          stdout << config.{{ field.id }}.url << '\n'
        else
          stdout << "not set".colorize.dark_gray << '\n'
        end

        stdout << "{{ field.id }} key: ".colorize.bold
        if config.{{ field.id }}?
          stdout << config.{{ field.id }}.key << '\n'
        else
          stdout << "not set".colorize.dark_gray << "\n\n"
        end
      {% end %}

      stdout << "ratelimit: ".colorize.bold << config.ratelimit << '\n'
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
            error "directory does not exist"
            return
          end

          path = Path[dir, ".soar.yml"]
        else
          path = path = Soar::Config::PATH
          Dir.mkdir_p path.dirname
        end

        if File.file?(path) && !options.has?("force")
          error "a config file already exists at this location"
          error "re-run with the '--force' flag to overwrite"
          return
        end

        begin
          File.write path, Soar::Config.new.to_yaml[4..]
        rescue ex
          error "failed to initialize config:"
          error ex
        end
      end
    end

    private class Copy < Base
      def setup : Nil
        @name = "copy"

        add_argument "src", required: true
        add_argument "dest", required: true
        add_option 'f', "force"
      end

      def run(arguments : Cling::Arguments, options : Cling::Options) : Nil
        src = arguments.get("src").as_s
        dest = arguments.get("dest").as_s
        return if src == dest

        case
        when src == "global"
          src = Soar::Config::PATH
          dest = Path[dest, ".soar.yml"] unless dest.ends_with? ".soar.yml"
        when dest == "global"
          dest = Soar::Config::PATH
          src = Path[src, ".soar.yml"] unless src.ends_with? ".soar.yml"
        else
          src = Path[src, ".soar.yml"] unless src.ends_with? ".soar.yml"
          dest = Path[dest, ".soar.yml"] unless dest.ends_with? ".soar.yml"
        end

        unless File.file? src
          error "source config not found"
          return
        end

        if File.file?(dest) && !options.has?("force")
          error "destination config already exists"
          return
        end

        File.copy src, dest
      end
    end
  end
end
