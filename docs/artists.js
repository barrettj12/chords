async function getArtists() {
  const resp = await fetch(SERVER_URL + '/artists')
  document.getElementById("location").innerHTML = await resp.text()
}