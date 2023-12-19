module Soar::Models
  class FractalList(T)
    include JSON::Serializable

    getter object : String
    getter data : Array(FractalItem(T))
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
    getter created_at : String
    getter updated_at : String?
  end
end
