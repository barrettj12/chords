let searchBar = document.getElementById("search-box");
let suggestionsHTML = document.getElementById("search-suggestions")

// Set up search
document.addEventListener('DOMContentLoaded', function() {
  searchBar.addEventListener("input", fillSuggestions);

  searchBar.addEventListener('focus', function() {
    suggestionsHTML.classList.add('show');
  });
  
  document.addEventListener('click', function(e) {
    if (!searchBar.contains(e.target) && !suggestionsHTML.contains(e.target)) {
      hideSuggestions();
    }
  });
});

// Hide search suggestions
function hideSuggestions() {
  suggestionsHTML.classList.remove('show');
}

async function fillSuggestions() {
  let query = searchBar.value
  let timestamp = new Date()

  if (!query) {
    // Search bar is empty, so exit early without calling API.
    suggestionsHTML.innerHTML = ''
    suggestionsHTML.setAttribute('data-last-updated', timestamp.toJSON())
    return
  }
  
  let resp = await fetch("/api/v0/search?" + new URLSearchParams({
    q: query,
  }))
  let suggestionsJSON = await resp.json()

  // It's possible that while we've been waiting, another request has
  // already been initiated, returned, and the suggestions updated.
  // If so, we don't want to overwrite suggestions from a later request.
  // Get the timestamp to check if this is the case.
  let lastUpdatedString = suggestionsHTML.getAttribute('data-last-updated')
  // For the first search, the 'data-last-updated' attribute is not
  // defined, but then the date is parsed as 1/1/1970, which still works.
  let lastUpdated = new Date(lastUpdatedString)
  if (lastUpdated > timestamp) {
    return
  }

  showSuggestions(suggestionsJSON)

  // suggestionsHTML.innerHTML = ''

  // if (suggestionsJSON.length == 0) {
  //   let row = suggestionsHTML.insertRow()
  //   let cell = row.insertCell()
  //   cell.setAttribute('style', "border: 1px solid black")
  //   cell.innerHTML = "No results found."
  // }

  // for (let res of suggestionsJSON) {
  //   let row = suggestionsHTML.insertRow()
  //   let cell = row.insertCell()
  //   cell.setAttribute('style', "border: 1px solid black")
    
  //   let link = document.createElement('a');
  //   if (res.type === "artist") {
  //     link.innerText = `${res.name} · artist`;
  //     link.setAttribute('href', "/b/songs?" + new URLSearchParams({
  //       artist: res.name,
  //     }));
  //   } else if (res.type === "song") {
  //     link.innerText = `${res.meta.name} · ${res.meta.artist} · ${res.meta.album}`;
  //     link.setAttribute('href', "/b/chords?" + new URLSearchParams({
  //       id: res.meta.id,
  //     }));
  //   }

  //   link.setAttribute('style', "text-decoration: none; color: black")
  //   cell.appendChild(link);
  // }

  suggestionsHTML.setAttribute('data-last-updated', timestamp.toJSON())
}

// Show search suggestions
function showSuggestions(results) {
  suggestionsHTML.innerHTML = '';

  if (results.length === 0) {
    const item = document.createElement('div');
    item.className = 'suggestion-item';
    item.textContent = 'No results found.';
    suggestionsHTML.appendChild(item);
    return;
  }
  
  results.forEach(result => {
    const item = document.createElement('div');
    item.className = 'suggestion-item';

    switch (result.type) {
      case "song":
        item.textContent = `${result.meta.name} · ${result.meta.artist}`;
        if (result.meta.album) {
          item.textContent += ` · ${result.meta.album}`;
        }
        item.onclick = () => {
          window.location.href = `/c/chords?id=${result.meta.id}`;
        };
        break;
        
      case "artist":
        item.textContent = `${result.name} · artist`;
        item.onclick = () => {
          window.location.href = `/c/songs?artist=${encodeURIComponent(result.name)}`;
        };
    }
    suggestionsHTML.appendChild(item);
  });
  
  suggestionsHTML.classList.add('show');
  // suggestionsHTML.setAttribute('data-last-updated', timestamp.toJSON())
}
