module Soar::Commands::App
  class Servers < Base
    def setup : Nil
      @name = "servers"

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
        servers, meta = request get: "/api/application/servers", as: Array(Models::Server)
        if servers.empty?
          unless options.has? "json"
            stdout << "Showing 0 results from page "
            stdout << meta.current_page << " of " << meta.total_pages << '\n'
            stdout << "\n┃ ".colorize.light_gray << "total: ".colorize.dark_gray << meta.total
            return
          end
        end

        if options.has? "json"
          servers.to_json stdout
          return
        end

        width = 2 + (Math.log(servers.last.id.to_f + 1) / Math.log(10)).ceil.to_i
        servers.each do |server|
          server.to_s(stdout, width)
          stdout << "\n\n"
        end

        stdout << "Showing " << meta.count << " results from page "
        stdout << meta.current_page << " of " << meta.total_pages << '\n'
        stdout << "\n┃ ".colorize.light_gray << "total: ".colorize.dark_gray << meta.total
      end
    end
  end
end
