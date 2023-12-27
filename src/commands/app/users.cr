module Soar::Commands::App
  class Users < Base
    def setup : Nil
      @name = "users"

      add_command List.new
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
      end

      def pre_run(arguments : Cling::Arguments, options : Cling::Options) : Nil
        super

        @filters << "username" if options.has? "username"
        @filters << "email" if options.has? "email"
        @filters << "uuid" if options.has? "uuid"
      end

      def run(arguments : Cling::Arguments, options : Cling::Options) : Nil
        path = "/api/application/users"

        unless @filters.empty?
          path += "?"
          path += URI::Params.build do |params|
            if options.has? "username"
              params.add "filter[username]", options.get("username").as_s
            end

            if options.has? "email"
              params.add "filter[email]", options.get("email").as_s
            end

            if options.has? "uuid"
              params.add "filter[uuid]", options.get("uuid").as_s
            end
          end
        end

        users, meta = request get: path, as: Array(Models::User)
        return if users.empty?

        width = 2 + (Math.log(users.last.id.to_f + 1) / Math.log(10)).ceil.to_i

        users.each do |user|
          Colorize.with.bold.on_light_gray.surround(stdout) do |io|
            io << user.id.to_s.center(width)
          end

          if external = user.external_id
            stdout << ' '
            Colorize.with.dark_gray.surround(stdout) do |io|
              io << '[' << external << ']'
            end
          end

          stdout << "\n\n┃ ".colorize.light_gray << "username:  ".colorize.bold
          stdout << user.username
          stdout << "\n┃ ".colorize.light_gray << "full name: ".colorize.bold
          stdout << user.first_name << ' ' << user.last_name

          stdout << "\n┃ ".colorize.light_gray << "email:     ".colorize.bold << user.email
          stdout << "\n┃ ".colorize.light_gray << "language:  ".colorize.bold << user.language

          stdout << "\n┃ ".colorize.light_gray << "is admin:  ".colorize.bold
          stdout << (user.root_admin? ? "true".colorize.green : "false".colorize.red)

          stdout << "\n┃ ".colorize.light_gray << "has 2FA:   ".colorize.bold
          stdout << (user.two_factor? ? "true".colorize.green : "false".colorize.red)

          stdout << "\n┃ ".colorize.light_gray << "created:   ".colorize.bold
          user.created_at.to_s(stdout, "%F %R")

          # TODO: implement 'x hours y minutes ago' format for updated at
          stdout << "\n┃ ".colorize.light_gray << "updated:   ".colorize.bold
          if updated = user.updated_at
            updated.to_s(stdout, "%F %R")
          else
            stdout << "N/A".colorize.dark_gray
          end

          stdout << "\n\n"
        end

        stdout << "Showing " << meta.count << " results from page "
        stdout << meta.current_page << " of " << meta.total_pages << '\n'
        stdout << "\n┃ ".colorize.light_gray << "total: ".colorize.dark_gray << meta.total
        stdout << "\n┃ ".colorize.light_gray

        if @filters.empty?
          stdout << "no filters applied"
        else
          stdout << "filters: ".colorize.dark_gray << @filters.join(", ")
        end
      end
    end
  end
end
