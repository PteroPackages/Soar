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
  end
end
