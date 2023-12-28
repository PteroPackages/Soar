module Soar::Models
  class FractalList(T)
    include JSON::Serializable

    getter object : String
    getter data : Array(FractalItem(T))
    @[JSON::Field(key: "meta", root: "pagination")]
    getter meta : FractalMeta
  end

  class FractalMeta
    include JSON::Serializable

    getter count : Int32
    getter total : Int32
    getter current_page : Int32
    getter per_page : Int32
    getter total_pages : Int32
  end

  class FractalItem(T)
    include JSON::Serializable

    getter object : String
    getter attributes : T
  end

  class User
    include JSON::Serializable

    getter id : Int32
    property external_id : String?
    getter uuid : String
    property username : String
    property email : String
    property first_name : String
    property last_name : String
    property language : String
    property? root_admin : Bool
    @[JSON::Field(key: "2fa")]
    property? two_factor : Bool
    getter created_at : Time
    getter updated_at : Time?

    def to_s(io : IO, width : Int32) : Nil
      Colorize.with.bold.on_light_gray.surround(io) do |_io|
        _io << id.to_s.center(width)
      end

      if external = external_id
        io << ' '
        Colorize.with.dark_gray.surround(io) do |_io|
          _io << '[' << external << ']'
        end
      end

      io << "\n\n┃ ".colorize.light_gray << "username:  ".colorize.bold
      io << username
      io << "\n┃ ".colorize.light_gray << "full name: ".colorize.bold
      io << first_name << ' ' << last_name

      io << "\n┃ ".colorize.light_gray << "email:     ".colorize.bold << email
      io << "\n┃ ".colorize.light_gray << "language:  ".colorize.bold << language

      io << "\n┃ ".colorize.light_gray << "is admin:  ".colorize.bold
      io << (root_admin? ? "true".colorize.green : "false".colorize.red)

      io << "\n┃ ".colorize.light_gray << "has 2FA:   ".colorize.bold
      io << (two_factor? ? "true".colorize.green : "false".colorize.red)

      io << "\n┃ ".colorize.light_gray << "created:   ".colorize.bold
      created_at.to_s(io, "%F %R")

      io << "\n┃ ".colorize.light_gray << "updated:   ".colorize.bold
      if updated = updated_at
        updated.to_s(io, "%F %R")
      else
        io << "N/A".colorize.dark_gray
      end
    end
  end

  class Server
    include JSON::Serializable

    getter id : Int32
    property external_id : String?
    getter uuid : String
    getter identifier : String
    property name : String
    property description : String?
    property status : Status = :none
    # getter limits
    # getter feature_limits
    property user : Int32
    getter node : Int32
    getter allocation : Int32
    getter nest : Int32
    getter egg : Int32
    # getter container
    getter created_at : Time
    getter updated_at : Time?

    enum Status
      None
      Installing
      InstallFailed
      Suspended
      RestoringBackup
    end

    def to_s(io : IO, width : Int32) : Nil
      base = case status
             when .none?
               :light_gray
             when .installing?, .restoring_backup?
               :yellow
             else
               :red
             end

      Colorize.with.bold.on(base).surround(io) do |_io|
        _io << id.to_s.center(width)
      end

      if external = external_id
        io << ' '
        Colorize.with.dark_gray.surround(io) do |_io|
          _io << '[' << external << ']'
        end
      end

      if !Colorize.enabled? && !status.none?
        io << '(' << status << ')'
      end

      io << '\n' << identifier << uuid.sub(identifier, "").colorize.dark_gray
      io << "\n\n┃ ".colorize.light_gray << "name:    ".colorize.bold << name

      io << "\n┃ ".colorize.light_gray << "about:   ".colorize.bold
      if (desc = description) && desc.presence
        if desc.size > 60
          desc = desc[0..60] + "...".colorize.dark_gray.to_s
        end

        io << desc
      else
        io << "N/A".colorize.dark_gray
      end

      io << "\n┃ ".colorize.light_gray << "status:  ".colorize.bold
      case status
      in .none?             then io << "N/A".colorize.dark_gray
      in .installing?       then io << "installing"
      in .install_failed?   then io << "install failed"
      in .suspended?        then io << "suspended"
      in .restoring_backup? then io << "restoring backup"
      end

      io << "\n┃ ".colorize.light_gray << "created: ".colorize.bold
      created_at.to_s(io, "%F %R")

      io << "\n┃ ".colorize.light_gray << "updated: ".colorize.bold
      if updated = updated_at
        updated.to_s(io, "%F %R")
      else
        io << "N/A".colorize.dark_gray
      end
    end
  end
end
