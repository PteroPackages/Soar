module Soar
  class Config
    include YAML::Serializable

    PATH = {% if flag?(:win32) %}
             Path[ENV["APPDATA"], "soar", ".soar.yml"]
           {% else %}
             Path[ENV["XDG_DATA_HOME"] || Path.home / ".local" / "share" / "soar", ".soar.yml"]
           {% end %}

    class Error < Exception
    end

    class Auth
      include YAML::Serializable

      property url : String = ""
      property key : String = ""

      def initialize
      end
    end

    class HTTPConfig
      include YAML::Serializable

      property? parse_body : Bool = true
      property? parse_errors : Bool = true
      property? parse_indent : Bool = true
      property? retry_ratelimit : Bool = true

      def initialize
      end
    end

    class LogConfig
      include YAML::Serializable

      property? use_color : Bool = true
      property? use_debug : Bool = false
      property? ignore_warnings : Bool = false

      def initialize
      end
    end

    @[YAML::Field(ignore: true)]
    property? resolved : Bool = false
    property! application : Auth
    property! client : Auth
    getter http : HTTPConfig
    getter logs : LogConfig

    def self.load : self
      load_local || load_global || Config.new
    end

    def self.load_local : self?
      load Path[Dir.current, ".soar.yml"]
    end

    def self.load_global : self?
      load PATH
    end

    def self.load(path : String | Path) : self?
      return nil unless File.file? path

      from_yaml(File.read path).tap &.resolved = true
    rescue ex : YAML::ParseException
      raise Error.new ex.message
    end

    def self.load_with_options(options : Cling::Options) : self
      config = load
      config.http.retry_ratelimit = true if options.has? "retry-ratelimit"
      config.http.retry_ratelimit = false if options.has? "no-retry-ratelimit"

      config.http.parse_body = true if options.has? "parse-body"
      config.http.parse_body = false if options.has? "no-parse-body"

      config.http.parse_errors = true if options.has? "parse-errors"
      config.http.parse_errors = false if options.has? "no-parse-errors"

      config.http.parse_indent = true if options.has? "parse-indent"
      config.http.parse_indent = false if options.has? "no-parse-indent"

      config
    end

    def initialize
      @application = Auth.new
      @client = Auth.new
      @http = HTTPConfig.new
      @logs = LogConfig.new
    end
  end
end
