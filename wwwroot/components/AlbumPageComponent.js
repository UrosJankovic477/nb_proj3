import {Component,h,createRef} from '/js/preact.js';
import htm from '/js/htm.js'
import { Album } from './AlbumComponent.js';
import {Song} from '/components/SongComponent.js';


const html = htm.bind(h)

class AlbumPage extends Component {
  albumartref = createRef("albumartref")
  fileref = createRef("fileref")
  nametxt = createRef("nametxt")
  pricetxt = createRef("pricetxt")
  genretxt = createRef("genretxt")

  async uploadImage(e) {
    let token = document.cookie.replace("token=", "");
    let img_file = this.albumartref.current.files[0]
    let album_uuid = this.state.albumInfo.uuid

    let new_form_data = new FormData()
    new_form_data.append("json_part", JSON.stringify({
      Token: token,
      AlbumUUID: album_uuid,
    }))
    new_form_data.append("art_bytes_part", img_file)
    let requestOptions = {
      method: 'POST',
      body: new_form_data,
      redirect: 'follow'
    };
    try {
      let response = await fetch("/api/uploadAlbumArt", requestOptions)
      window.location.reload()
    } 
    catch (Fehler) {
        return Fehler
    }

  }

  renderUploadButton() {
    if (!this.state.albumInfo.author) {
        return null
    }
    return html `<div>
    <input ref=${this.albumartref} type="file"/>
    <button onClick=${(e) => this.uploadImage(e)}>Upload Image</button>
    </div>`
  }


    async componentDidMount() {
      let album_uuid = window.location.pathname.replace("/album/", "")
      var myHeaders = new Headers();
      myHeaders.append("Content-Type", "application/json");
      let token = document.cookie.replace("token=", "");
      var raw = JSON.stringify({
        Token: token,
        AlbumUUID: album_uuid,
      });
      var requestOptions = {
        method: 'POST',
        headers: myHeaders,
        body: raw,
        redirect: 'follow'
      };
      try {
          let response = await fetch("/api/album", requestOptions)
          let result = await response.json()
          let genre_set = new Set()
          let song_list = []
          result.Songs.forEach(s => {
            genre_set.add(s.Genre)
            song_list.push({
              songname: s.Title,
              songduration: s.Length,
              price: s.Price,
              uuid: s.UUID,
              genre: s.Genre
            })
          })
    
          let albumInfo = {
            uuid: result.UUID,
            title: result.Title,
            price: result.Price,
            songList: song_list,
            genres: genre_set,
            author: result.Created,
            owns: result.Owns,
          }
          this.setState({albumInfo:albumInfo})
             
      } catch (Fehler) {
          return Fehler
      }     
    }
   
    async uploadSong(e) {
      var myHeaders = new Headers();
      let token = document.cookie.replace("token=", "");
      let album_uuid = this.state.albumInfo.uuid
      let title = this.nametxt.current.value
      let price = Number(this.pricetxt.current.value)
      let genre = this.genretxt.current.value
      let song_bytes =  this.fileref.current.files[0]
    
      let new_form_data = new FormData()

      let raw = {
        Token: token,
        AlbumUUID: album_uuid,
        UUID: "",
        Title: title,
        Price: price,
        Genre: genre,
      }
      
      new_form_data.append("json_part", JSON.stringify(raw))
      new_form_data.append("song_bytes_part", song_bytes)

      let requestOptions = {
        method: 'POST',
        body: new_form_data,
        redirect: 'follow'
      };
      try {
        let response = await fetch("/api/addSongToAlbum", requestOptions)
        window.location.reload()
      } 
      catch (Fehler) {
          return Fehler
      }

    }
    renderSong() {
      let htmllist = []
      if(!this.state.albumInfo) {return null}
      const songlist = this.state.albumInfo.songList
      let ctr = 1

      songlist.forEach((song) => {
        const songduration = new Date(song.songduration * 1000).toISOString().slice(11, 19);
        const price = song.price.toFixed(2)
        if (this.state.albumInfo.author) {
          htmllist.push( html `<${Song} num="${ctr}" songname="${song.songname}" createdSong  songlength="${songduration}" price="${price}" uuid="${song.uuid}" />`)
        }
        else if(this.state.albumInfo.owns) {
          htmllist.push( html `<${Song} num="${ctr}" songname="${song.songname}" ownsSong  songlength="${songduration}" price="${price}" uuid="${song.uuid}" />`)
        }
        else {
          htmllist.push( html `<${Song} num="${ctr}" songname="${song.songname}"  songlength="${songduration}" price="${price}" uuid="${song.uuid}" />`)
        }
        ctr++
      })
      if(this.state.albumInfo.author) {
        htmllist.push(html`<div class="songUploadDiv"><input ref=${this.fileref} type="file"/> <input ref=${this.nametxt} placeholder="Song name"/><label for="pricetxt">$</label><input type="number" id="pricetxt" ref=${this.pricetxt}/><input ref=${this.genretxt} placeholder="Genre"/><button onClick=${(e) => this.uploadSong(e)}>Add Song</button></div>`)
      }
      

      return htmllist
    }

    renderAlbumInfo() {
      let htmllist = []
      const albumInfo = this.state.albumInfo

      if (!albumInfo) {
        return null
      }


      albumInfo.genres.forEach((genre) => {
        htmllist.push(html `<p class=albumGenre>${genre}</p>`)
       
      })

      return html `<div class="basicAlbumInfo albumInfo">
        <img src="/albumart/${albumInfo.uuid}"></img>
        ${this.renderUploadButton()}
        <p>${albumInfo.albumName}</p>
        <p>${albumInfo.albumartistdisplayName}</p>
        <p>$${albumInfo.price.toFixed(2)}</p>
        <h1>Genres</h1>
        ${htmllist}
        </div>`
    }
   

    render() {
      return html `
      <div id="background" class="albumBackground" style="background-image:url(/albumart/${ this.state.albumInfo ? this.state.albumInfo.uuid : ""}"></div><div id="albumView"><div class="smallerDiv" id="leftAlbumDiv">${this.renderAlbumInfo()}</div><div class="biggerDiv" id="rightAlbumDiv">${this.renderSong()}</div></div>
      `
      

    }
}
export {AlbumPage}
