module Soar::Commands::App
  class Locations < Base
    def setup : Nil
      @name = "locations"

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

        add_option "short", type: :single
        add_option "long", type: :single
        add_option "json"
        add_option "page", type: :single
        add_option "per-page", type: :single
        add_option "sort", type: :single
      end

      def pre_run(arguments : Cling::Arguments, options : Cling::Options) : Bool
        return false unless super

        @filters << "short" if options.has? "short"
        @filters << "long" if options.has? "long"

        true
      end

      def run(arguments : Cling::Arguments, options : Cling::Options) : Nil
        path = "/api/application/locations?"
        def_base_and_filter_params short: "short", long: "long"

        locations, meta = request get: path, as: Array(Models::Location)
        if options.has? "json"
          locations.to_json stdout
          return
        end

        unless locations.empty?
          width = 2 + (Math.log(locations.last.id.to_f + 1) / Math.log(10)).ceil.to_i
          locations.each do |location|
            location.to_s(stdout, width)
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
        add_option "json"
      end

      def run(arguments : Cling::Arguments, options : Cling::Options) : Nil
        location = request get: "/api/application/locations/#{arguments.get("id")}", as: Models::Location

        if options.has? "json"
          location.to_json stdout
          return
        end

        width = 2 + (Math.log(location.id.to_f + 1) / Math.log(10)).ceil.to_i
        location.to_s(stdout, width)
        stdout.puts
      end
    end

    private class Create < Base
      def setup : Nil
        @name = "create"

        add_option 'd', "data", type: :single
        add_option 'i', "input"
      end

      def pre_run(arguments : Cling::Arguments, options : Cling::Options) : Bool
        return false unless super

        has_data = options.has? "data"
        has_input = options.has? "input"

        if has_data && has_input
          error "cannot specify 'data' and 'input' flag; pick one"
          return false
        elsif !has_data && !has_input
          error "either 'data' or 'input' option must be specified"
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
          input = Resolver.parse_json_or_map options.get("data").as_s
        end

        location = request post: "/api/application/locations", data: input, as: Models::Location
        width = 2 + (Math.log(location.id.to_f + 1) / Math.log(10)).ceil.to_i
        location.to_s(stdout, width)
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
        request delete: "/api/application/locations/#{id}"
      end
    end
  end
end
