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
        users = request get: "/api/application/users", as: Array(Models::User)

        stdout << users.size << " Result" << (users.size == 1 ? '\n' : "s\n")
        return if users.empty?
        stdout.puts

        users.each do |user|
          stdout << ' ' << user.id.colorize.bold.on_light_gray << ' ' << user.username
          if external = user.external_id
            stdout << " (" << external << ")"
          end

          stdout << "\n " << user.first_name << ' ' << user.last_name
          stdout << " <" << user.email << ">\n\n"
        end
      end
    end
  end
end
