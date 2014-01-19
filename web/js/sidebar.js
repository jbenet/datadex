(function ($) {

  $(function(){

    // IE10 viewport hack for Surface/desktop Windows 8 bug
    //
    // See Getting Started docs for more information
    if (navigator.userAgent.match(/IEMobile\/10\.0/)) {
      var msViewportStyle = document.createElement("style");
      msViewportStyle.appendChild(
        document.createTextNode(
          "@-ms-viewport{width:auto!important}"
        )
      );
      document.getElementsByTagName("head")[0].
        appendChild(msViewportStyle);
    }


    var $window = $(window)
    var $body   = $(document.body)

    var navHeight = $('.navbar').outerHeight(true) + 10

    $body.scrollspy({
      target: '.bs-sidebar',
      offset: navHeight
    })

    $window.on('load', function () {
      $body.scrollspy('refresh')
    })

    var navOffset = function() {
      var $sideBar = $('.bs-sidebar')
      var sideBarMargin  = parseInt($sideBar.children(0).css('margin-top'), 10)
      var navOuterHeight = $('.datadex-nav').height()
      return navOuterHeight + sideBarMargin
    }

    var setNavHeight = function() {
      var wh = $(window).height()
      var h = wh - navOffset()
      if ($(window).width() > 768 && wh < $('.bs-sidenav').height()) {
        $('.bs-sidenav').css("max-height", h)
      } else {
        $('.bs-sidenav').css("max-height", "")
      }
    };

    // back to top
    setTimeout(function () {
      var $sideBar = $('.bs-sidebar')

      $sideBar.affix({
        offset: {
          top: function () {
            setNavHeight()
            var offsetTop = $sideBar.offset().top;
            return (this.top = offsetTop - navOffset() - 0)
          }
        , bottom: function () {
            setNavHeight()
            return (this.bottom = $('footer').outerHeight(true) + 20)
          }
        }
      })

      setNavHeight()

      // find current page links and color them.
      $sideBar.find("[href]").each(function() {
        if (this.href == window.location.href ||
            this.href == window.location.pathname ||
            this.href == window.location.pathname + window.location.hash) {
          $(this).addClass("current-page");
        }
      });

    }, 100)

})

})(jQuery);
