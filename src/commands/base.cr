module Soar::Commands
  abstract class Base < Cling::Command
    def initialize
      super

      @inherit_options = true
      add_option "no-color"
      add_option 'h', "help"
    end

    def pre_run(arguments : Cling::Arguments, options : Cling::Options) : Bool
      Colorize.enabled = false if options.has? "no-color"

      if options.has? "help"
        stdout.puts help_template

        false
      else
        true
      end
    end

    def on_error(ex)
      pp! ex
    end
  end
end
