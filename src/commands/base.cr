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

    def on_error(ex : Exception)
      case ex
      when Cling::CommandError
        error ex.to_s
        error "see 'soar --help' for more information"
      else
        error "unexpected exception:"
        error ex.to_s
        error "please report this on the Soar GitHub issues page:"
        error "https://github.com/PteroPackages/Soar/issues"
      end
    end

    def on_missing_arguments(args : Array(String))
      error "missing required argument#{"s" if args.size > 1} for this command:"
      error args.join ", "
    end

    def on_unknown_arguments(args : Array(String))
      error "unexpected argument#{"s" if args.size > 1} for this command:"
      error args.join ", "
    end

    def on_invalid_option(message : String)
      error message
    end

    def on_unknown_options(options : Array(String))
      error "unexpected option#{"s" if options.size > 1} for this command:"
      error options.join ", "
    end

    protected def warn(message : String) : Nil
      stdout << "warn".colorize.yellow << ": " << message << '\n'
    end

    protected def error(message : String) : Nil
      stderr << "error".colorize.red << ": " << message << '\n'
    end

    # TODO: show backtrace if debug mode
    protected def error(ex : Exception) : Nil
      error ex.to_s
    end
  end
end
