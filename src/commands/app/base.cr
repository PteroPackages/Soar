module Soar::Commands::App
  abstract class Base < Soar::Commands::Base
    protected getter config : Soar::Config { raise "unreachable" }
    protected getter client : Crest::Resource { raise "unreachable" }

    def pre_run(arguments : Cling::Arguments, options : Cling::Options) : Bool?
      super

      @config = Soar::Config.load_with_options options
      @client = Crest::Resource.new(
        config.application.url,
        headers: {
          "Authorization" => "Bearer #{config.application.key}",
          "Content-Type"  => "application/vnd.pterodactyl.v1+json",
          "Accept"        => "application/vnd.pterodactyl.v1+json",
        })

      Colorize.enabled = false unless config.logs.use_color?
    end
  end
end
