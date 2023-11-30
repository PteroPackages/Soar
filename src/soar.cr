require "cling"
require "crest"
require "yaml"

require "./config"

module Soar
  VERSION = "0.2.0"

  class CLI < Cling::Command
    def setup : Nil
      @name = "soar"
    end

    def run(arguments : Cling::Arguments, options : Cling::Options) : Nil
      stdout.puts help_template
    end
  end
end
