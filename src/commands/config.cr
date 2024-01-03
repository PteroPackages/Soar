module Soar::Commands
  class Config < Base
    def setup : Nil
      @name = "config"

      add_command Init.new
      add_command Set.new
      add_command Copy.new

      add_option 'g', "global"
      add_option 'l', "local"
    end

    def pre_run(arguments : Cling::Arguments, options : Cling::Options) : Bool
      return false unless super

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

      if config.nil?
        error "failed to load #{global ? "global" : "local"} config"
        system_exit
      end

      stdout << "app url: ".colorize.bold
      if config.app?
        stdout << config.app.url << '\n'
      else
        stdout << "not set".colorize.dark_gray << '\n'
      end

      stdout << "app key: ".colorize.bold
      if config.app?
        stdout << config.app.key << "\n\n"
      else
        stdout << "not set".colorize.dark_gray << "\n\n"
      end

      stdout << "client url: ".colorize.bold
      if config.client?
        stdout << config.client.url << '\n'
      else
        stdout << "not set".colorize.dark_gray << '\n'
      end

      stdout << "client key: ".colorize.bold
      if config.client?
        stdout << config.client.key << "\n\n"
      else
        stdout << "not set".colorize.dark_gray << "\n\n"
      end

      stdout << "ratelimit: ".colorize.bold << config.ratelimit << '\n'
    end

    private class Init < Base
      def setup : Nil
        @name = "init"

        add_argument "dir"
        add_option 'f', "force"
      end

      def run(arguments : Cling::Arguments, options : Cling::Options) : Nil
        if dir = arguments.get?("dir").try &.as_s
          unless Dir.exists? dir
            error "directory does not exist"
            system_exit
          end

          path = Path[dir, ".soar.yml"]
        else
          path = Soar::Config::PATH
          Dir.mkdir_p path.dirname
        end

        if File.file?(path) && !options.has?("force")
          error "a config file already exists at this location"
          error "re-run with the '--force' flag to overwrite"
          system_exit
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

        if src == "global"
          src = Soar::Config::PATH
          dest = Path[dest, ".soar.yml"] unless dest.ends_with? ".soar.yml"
        elsif dest == "global"
          dest = Soar::Config::PATH
          src = Path[src, ".soar.yml"] unless src.ends_with? ".soar.yml"
        else
          src = Path[src, ".soar.yml"] unless src.ends_with? ".soar.yml"
          dest = Path[dest, ".soar.yml"] unless dest.ends_with? ".soar.yml"
        end

        unless File.file? src
          error "source config not found"
          system_exit
        end

        if File.file?(dest) && !options.has?("force")
          error "destination config already exists"
          system_exit
        end

        File.copy src, dest
      end
    end

    private class Set < Base
      def setup : Nil
        @name = "set"

        add_argument "key"
        add_argument "value"
        add_option 'i', "input"
        add_option 'g', "global"
        add_option 'l', "local"
      end

      def pre_run(arguments : Cling::Arguments, options : Cling::Options) : Bool
        return false unless super

        unless options.has? "input"
          args = [] of String
          args << "key" unless arguments.has? "key"
          args << "value" unless arguments.has? "value"
          return true if args.empty?

          on_missing_arguments args
          return false
        end

        true
      end

      def run(arguments : Cling::Arguments, options : Cling::Options) : Nil
        if options.has? "input"
          if stdin.closed?
            error "cannot read from input file (already closed)"
            system_exit
          end

          input = Resolver.parse_json_or_map stdin.gets_to_end.chomp
        else
          input = {arguments.get("key").as_s => arguments.get("value").as_s}
        end

        global = options.has? "global"
        local = options.has? "local"
        config = if global
                   Soar::Config.load_global
                 elsif local
                   Soar::Config.load_local
                 else
                   Soar::Config.load
                 end

        if config.nil?
          error "failed to load #{global ? "global" : "local"} config"
          system_exit
        end

        input.each do |key, value|
          case key
          when "app.url"
            config.app.url = value
          when "app.key"
            config.app.key = value
          when "client.url"
            config.client.url = value
          when "client.key"
            config.client.key = value
          when "ratelimit"
            config.ratelimit = value
          else
            warn "invalid config key '#{key}'"
          end
        end

        path = global ? Soar::Config::PATH : Path[Dir.current, ".soar.yml"]
        File.write path, config.to_yaml[4..]
      end
    end
  end
end
