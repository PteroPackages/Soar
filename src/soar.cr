require "cling"
require "colorize"
require "crest"
require "json"
require "yaml"

require "./commands/*"
require "./commands/app/main"
require "./config"
require "./models"
require "./resolver"

module Soar
  VERSION    = "0.2.0"
  BUILD_DATE = {% if flag?(:win32) %}
                 {{ `powershell.exe -NoProfile Get-Date -Format "yyyy-MM-dd"`.stringify.chomp }}
               {% else %}
                 {{ `date +%F`.stringify.chomp }}
               {% end %}
  BUILD_HASH = {{ env("BUILD_HASH") || `git rev-parse HEAD`.stringify[0...8] }}

  class CLI < Commands::Base
    def setup : Nil
      @name = "soar"

      add_command Commands::App::Main.new
      add_command Commands::Config.new
      add_command Commands::Version.new
    end

    def run(arguments : Cling::Arguments, options : Cling::Options) : Nil
      stdout.puts help_template
    end
  end

  class SystemExit < Exception
  end
end
