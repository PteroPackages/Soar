module Soar::Commands::App
  class Locations < Base
    def setup : Nil
      @name = "locations"

      add_command List.new
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
      end

      def pre_run(arguments : Cling::Arguments, options : Cling::Options) : Bool
        return false unless super

        @filters << "short" if options.has? "short"
        @filters << "long" if options.has? "long"

        true
      end

      def run(arguments : Cling::Arguments, options : Cling::Options) : Nil
        path = "/api/application/locations"
        unless @filters.empty?
          def_filter_params short: "short", long: "long"
        end

        locations, meta = request get: path, as: Array(Models::Location)
        if locations.empty?
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
          locations.to_json stdout
          return
        end

        width = 2 + (Math.log(locations.last.id.to_f + 1) / Math.log(10)).ceil.to_i
        locations.each do |location|
          location.to_s(stdout, width)
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
        stdout.puts
      end
    end
  end
end
