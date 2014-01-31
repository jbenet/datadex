class datadex.Autocomplete

  constructor: (field, @options={}) ->
    @field = $(field)

    @field.typeahead
      name: @options.name || @field.attr('id') || 'search'
      template: @template
      remote:
        url: @url()
        filter: @filter

    @field.on 'typeahead:autocompleted', @selected
    @field.on 'typeahead:selected', @selected


  url: =>
    es_url = datadex.config.ELASTICSEARCH_URL
    source = 'source={"query":{"match_phrase_prefix":{"_all":"%QUERY"}}}'
    if @options.type
      "#{es_url}/#{@options.type}/_search/?#{source}"
    else
      "#{es_url}/_search/?#{source}"


  template: _.template '''
    <span class="icon">
      <i class="<%= type_icon %>"></i>
    </span>
    <span class="value"><%= value %></span>
    <span class="tagline"><%= tagline %></span>
    '''


  filter: (data) =>
    # deferring is a hack, but there is currently no event to use.
    # See https://github.com/twitter/typeahead.js/issues/130
    _.defer =>
      console.log('deferred MathJax typesetting');
      # MathJax.Hub.Queue(["Typeset", MathJax.Hub]);

    icons =
      User: 'icon-user'
      Dataset: 'icon-file'

    objects = _.map data.hits.hits, (object) =>
      key = datadex.util.key(object._source.key)
      _.extend {}, object._source,
        keyfn: key
        value: object._source.name
        tokens: object._source.name
        tagline: object._source.tagline
        type_icon: icons[key.type]

    _.filter objects, (object) =>
      object.keyfn.type in ['User', 'Dataset']



  selected: (event, object, dataset) =>


class datadex.Search extends datadex.Autocomplete
  selected: (event, object, dataset) =>
    document.location.href = object.link
