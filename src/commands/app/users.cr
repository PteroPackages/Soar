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
        #
      end
    end
  end
end
