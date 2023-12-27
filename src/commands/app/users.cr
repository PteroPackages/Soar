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
      def setup : Nil
        @name = "list"
      end

      def run(arguments : Cling::Arguments, options : Cling::Options) : Nil
        users, meta = request get: "/api/application/users", as: Array(Models::User)
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
        stdout << "\n┃ ".colorize.light_gray << "no filters or sorts applied".colorize.dark_gray
      end
    end
  end
end
