module Soar::Commands::App
  class Servers < Base
    def setup : Nil
      @name = "servers"

      add_command List.new
      add_command Get.new
      add_command Create.new
      add_command Delete.new
    end

    def run(arguments : Cling::Arguments, options : Cling::Options) : Nil
      stdout.puts help_template
    end

    private class List < Base
      @filters = [] of String

      def setup : Nil
        @name = "list"

        add_option "name", type: :single
        add_option "desc", type: :single
        add_option "uuid", type: :single
        add_option "image", type: :single
        add_option "json"
        add_option "page", type: :single
        add_option "per-page", type: :single
        add_option "sort", type: :single
      end

      def pre_run(arguments : Cling::Arguments, options : Cling::Options) : Bool
        return false unless super

        @filters << "name" if options.has? "name"
        @filters << "description" if options.has? "desc"
        @filters << "uuid" if options.has? "uuid"
        @filters << "image" if options.has? "image"

        true
      end

      def run(arguments : Cling::Arguments, options : Cling::Options) : Nil
        path = "/api/application/servers?"
        def_base_and_filter_params name: "name", desc: "description", uuid: "uuid", image: "image"

        servers, meta = request get: path, as: Array(Models::Server)
        if options.has? "json"
          servers.to_json stdout
          return
        end

        unless servers.empty?
          width = 2 + (Math.log(servers.last.id.to_f + 1) / Math.log(10)).ceil.to_i
          servers.each do |server|
            server.to_s(stdout, width)
            stdout << "\n\n"
          end
        end

        meta.to_s(stdout, @filters, options.get?("sort").try &.as_s)
      end
    end

    private class Get < Base
      def setup : Nil
        @name = "get"

        add_argument "id", required: true
        add_option 'e', "external"
        add_option "json"
      end

      def run(arguments : Cling::Arguments, options : Cling::Options) : Nil
        arg = arguments.get("id").as_s

        if options.has? "external"
          server = request get: "/api/application/servers/external/#{arg}", as: Models::Server
        else
          server = request get: "/api/application/servers/#{arg}", as: Models::Server
        end

        if options.has? "json"
          server.to_json stdout
          return
        end

        width = 2 + (Math.log(server.id.to_f + 1) / Math.log(10)).ceil.to_i
        server.to_s(stdout, width)
        stdout.puts
      end
    end

    private class Create < Base
      def setup : Nil
        @name = "create"

        add_argument "data", required: true
        add_option 'i', "input"
      end

      def pre_run(arguments : Cling::Arguments, options : Cling::Options) : Bool
        return false unless super

        unless arguments.has?("data") || options.has?("input")
          on_missing_arguments %w[data]
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

          input = stdin.gets_to_end.chomp
        else
          input = arguments.get("data").as_s
        end

        server = request post: "/api/application/servers", data: input, as: Models::Server
        width = 2 + (Math.log(server.id.to_f + 1) / Math.log(10)).ceil.to_i
        server.to_s(stdout, width)
        stdout.puts
      end
    end

    private class Delete < Base
      def setup : Nil
        @name = "delete"

        add_argument "id", required: true
      end

      def run(arguments : Cling::Arguments, options : Cling::Options) : Nil
        id = arguments.get("id").as_s
        request delete: "/api/application/servers/#{id}"
      end
    end
  end
end
