module Soar::Commands
  class Version < Base
    def setup : Nil
      @name = "version"
    end

    def run(arguments : Cling::Arguments, options : Cling::Options) : Nil
      stdout << "soar version " << Soar::VERSION << '\n'
      # stdout << " [" << Soar::BUILD_HASH << "] ("
      # stdout << Soar::BUILD_DATE << ")\n"
    end
  end
end
