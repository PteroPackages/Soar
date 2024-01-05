module Soar
  class Config
    include YAML::Serializable

    PATH = {% if flag?(:win32) %}
             Path[ENV["APPDATA"], "soar", ".soar.yml"]
           {% else %}
             Path[ENV["XDG_DATA_HOME"]? || Path.home / ".local" / "share" / "soar", ".soar.yml"]
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

    @[YAML::Field(ignore: true)]
    property? resolved : Bool = false
    property! app : Auth
    property! client : Auth
    property! ratelimit : String

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

    # TODO: will this still be needed?
    def self.load_with(options : Cling::Options) : self
      config = load
      if value = options.get?("cfg-url")
        config.app.url = config.client.url = value.as_s
      end

      if value = options.get?("cfg-key")
        config.app.key = config.client.key = value.as_s
      end

      if value = options.get?("retry-ratelimit")
        config.ratelimit = value.as_s
      end

      config
    end

    def initialize
      @app = Auth.new
      @client = Auth.new
    end
  end
end
