class datadex.util

  @urlName: (name) =>
    name.trim().replace(/\s+/gi, '-')


  @post: (url, data, callback) =>

    $.ajax
      url: url
      data: JSON.stringify(data)
      method: 'POST'
      contentType: 'application/json; charset=utf-8'
      dataType: 'json'
      success: (data) => callback data
      failure: (err) => console.log err

  @key: (keystr) =>
    key = => String(keystr)
    key.base = _.last key().split("/")
    key.name = key.base.split(":")[-1] || ""
    key.type = key.base.split(":")[0] || ""
    key
