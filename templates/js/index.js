function submitSearch() {
    $.ajax({
      url: "/search",
      method: "POST",
      data: $("#search-form").serialize(),
      success: function(rawData) {
        var parsed = JSON.parse(rawData);
        if (!parsed) return;

        var searchResults = $("#search-results");
        searchResults.empty();

        parsed.forEach(function(result) {
          var row = $("<tr><td>" + result.Title + "</td><td>" + result.Author + "</td><td>" + result.Year +  "</td><td>" + result.ID + "</td></tr>");
          searchResults.append(row);
        });
      }
    });

    return false;
  }