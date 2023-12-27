module Soar::Resolver
  def self.parse_json_or_map(source : String) : Hash(String, String)
    if source.starts_with? '{'
      return JSON.parse(source).as_h.transform_values &.as_s
    end

    reader = Char::Reader.new source
    data = {} of String => String
    on_value = false
    key = ""

    loop do
      case reader.current_char
      when '\0'
        raise "Missing value pair for key" if on_value
        break
      when ' '
        reader.next_char
      when '='
        if reader.next_char == '\0'
          raise "Missing value pair for key"
        end

        if reader.current_char == '"'
          reader.next_char
          start = reader.pos

          loop do
            case reader.current_char
            when '\0' then raise "Unterminated quote string for value pair"
            when '"'  then break
            else           reader.next_char
            end
          end

          reader.next_char
          data[key] = reader.string[start...reader.pos - 1]
          on_value = false
        else
          start = reader.pos
          while reader.has_next? && reader.current_char != ' '
            reader.next_char
          end

          data[key] = reader.string[start...reader.pos]
          on_value = false
        end
      when .ascii_alphanumeric?
        start = reader.pos
        while reader.has_next? && reader.current_char != '='
          reader.next_char
        end

        key = reader.string[start...reader.pos]
        on_value = true
      end
    end

    data
  end
end
