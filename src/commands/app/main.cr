require "./base"
require "./servers"
require "./users"

module Soar::Commands::App
  class Main < Base
    def setup : Nil
      @name = "app"

      add_command App::Servers.new
      add_command App::Users.new
    end

    def run(arguments : Cling::Arguments, options : Cling::Options) : Nil
      stdout.puts help_template
    end
  end
end
