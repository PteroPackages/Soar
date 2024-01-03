module Soar::Commands::App
  class Users < Base
    def setup : Nil
      @name = "users"

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

        add_option "username", type: :single
        add_option "email", type: :single
        add_option "uuid", type: :single
        add_option "json"
      end

      def pre_run(arguments : Cling::Arguments, options : Cling::Options) : Bool
        return false unless super

        @filters << "username" if options.has? "username"
        @filters << "email" if options.has? "email"
        @filters << "uuid" if options.has? "uuid"

        true
      end

      def run(arguments : Cling::Arguments, options : Cling::Options) : Nil
        path = "/api/application/users"
        unless @filters.empty?
          def_filter_params username: "username", email: "email", uuid: "uuid"
        end

        users, meta = request get: path, as: Array(Models::User)
        if users.empty?
          unless options.has? "json"
            stdout << "Showing 0 results from page "
            stdout << meta.current_page << " of " << meta.total_pages << '\n'
            stdout << "\n┃ ".colorize.light_gray << "total: ".colorize.dark_gray << meta.total
            stdout << "\n┃ ".colorize.light_gray

            if @filters.empty?
              stdout << "no filters applied"
            else
              stdout << "filters: ".colorize.dark_gray << @filters.join(", ")
            end
            stdout.puts

            return
          end
        end

        if options.has? "json"
          users.to_json stdout
          return
        end

        width = 2 + (Math.log(users.last.id.to_f + 1) / Math.log(10)).ceil.to_i
        users.each do |user|
          user.to_s(stdout, width)
          stdout << "\n\n"
        end

        stdout << "Showing " << meta.count << " results from page "
        stdout << meta.current_page << " of " << meta.total_pages << '\n'
        stdout << "\n┃ ".colorize.light_gray << "total: ".colorize.dark_gray << meta.total
        stdout << "\n┃ ".colorize.light_gray

        if @filters.empty?
          stdout << "no filters applied".colorize.dark_gray
        else
          stdout << "filters: ".colorize.dark_gray << @filters.join(", ")
        end
        stdout.puts
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
          user = request get: "/api/application/users/external/#{arg}", as: Models::User
        else
          user = request get: "/api/application/users/#{arg}", as: Models::User
        end

        if options.has? "json"
          user.to_json stdout
          return
        end

        width = 2 + (Math.log(user.id.to_f + 1) / Math.log(10)).ceil.to_i
        user.to_s(stdout, width)
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
            return
          end

          input = Resolver.parse_json_or_map stdin.gets_to_end.chomp
        else
          input = Resolver.parse_json_or_map options.get("data").as_s
        end

        {"username", "email", "first_name", "last_name"}.each do |name|
          unless input.has_key? name
            error "Missing required key '#{name}' to create user"
            return
          end
        end

        user = request post: "/api/application/users", data: input, as: Models::User
        width = 2 + (Math.log(user.id.to_f + 1) / Math.log(10)).ceil.to_i
        user.to_s(stdout, width)
        stdout.puts
      end
    end

    private class Delete < Base
      def setup : Nil
        @name = "delete"

        add_argument "id", required: true
      end

      def run(arguments : Cling::Arguments, options : Cling::Options) : Nil
        id = arguments.get("id").as_s # TODO: parse to integer once Cling upstream is fixed
        request delete: "/api/application/users/#{id}"
      end
    end
  end
end
