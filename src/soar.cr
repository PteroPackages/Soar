require "cling"
require "crest"
require "yaml"

require "./commands/*"
require "./config"

module Soar
  VERSION = "0.2.0"

  class CLI < Commands::Base
    def setup : Nil
      @name = "soar"

      add_command Commands::Version.new
    end

    def run(arguments : Cling::Arguments, options : Cling::Options) : Nil
      stdout.puts help_template
    end
  end
end
