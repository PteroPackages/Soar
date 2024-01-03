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
      when SystemExit
        exit 1
      when Cling::CommandError
        error ex.to_s
        error "see 'soar --help' for more information"
      when Crest::RequestFailed
        case ex.http_code
        when .in?(400..499)
          data = Models::FractalError.from_json(ex.response.body).errors

          error "#{data.size} error#{"s" if data.size > 1} received"
          data.each do |err|
            stderr << "\n┃ ".colorize.bold << '[' << err.status << "] "
            stderr << err.code.colorize.bold << '\n'
            stderr << "┃ ".colorize.bold << err.detail << '\n'
          end
        when .in?(500..512)
          error "unexpected response: #{ex.http_code}"
          error "data:"
          stderr.puts ex.response.body
        else
          error "unknown http status: #{ex.http_code}"
          error "request cancelled"
        end
      else
        error "unexpected exception:"
        error ex.to_s
        error "please report this on the Soar GitHub issues page:"
        error "https://github.com/PteroPackages/Soar/issues"
      end

      exit 1
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

    protected def system_exit : NoReturn
      raise SystemExit.new
    end
  end
end
