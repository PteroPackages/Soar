module Soar::Commands
  class Config < Base
    def setup : Nil
      @name = "config"
    end

    def run(arguments : Cling::Arguments, options : Cling::Options) : Nil
      config = Soar::Config.load
      puts config.to_yaml[4..]
    end
  end
end
