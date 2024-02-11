import {Component,h,createRef} from '/js/preact.js';
import htm from '/js/htm.js'


const html = htm.bind(h)

class Song extends Component {
    audio = createRef("audio")
    constructor(props) {

       
        super(props);
        this.state = {playing:false}
        
        
      }

    togglePlay(e) {
    const player = this.audio.current
    let playinglocal = this.state.playing
    playinglocal = !playinglocal
    playinglocal ? player.play() : player.pause()
        this.setState({playing:playinglocal})

    }


    async deleteSong(e) {
      var myHeaders = new Headers();
            myHeaders.append("Content-Type", "application/json");
            let token = document.cookie.replace("token=", "");
            let song_uuid = this.props.uuid
            

            var raw = JSON.stringify({
              Token: token,
              SongUUID: song_uuid,
            });

            var requestOptions = {
              method: 'POST',
              headers: myHeaders,
              body: raw,
              redirect: 'follow'
            };
            try {
                let response = await fetch("/api/removeSongFromAlbum", requestOptions)
                window.location.reload()
                
            } catch (Fehler) {
                return Fehler
            }
    }

    addToCart(e) {

    }
    renderDeleteSong() {
      if(this.props.createdSong === undefined) {
        return null
      }
      return html`<p class="textbutton" onClick=${(e) => {this.deleteSong(e)}}>❌</p>`
    }

    renderDownloadSong() {
    if(!(this.props.ownsSong !== undefined || this.props.createdSong !== undefined)) {
        return null
      }
      return html`<p class="textbutton" onClick=${(e) => window.location.pathname = `/download/${this.props.uuid}`}>⬇️</p>`

    }

    render() {
      return html `
      <div class="songDiv">
      <p class="textbutton" onClick=${(e) => this.togglePlay(e)}>${this.state.playing ? '⏸️' : '▶️'}</p>
      ${this.renderDownloadSong()}
      <p class="songNum">${this.props.num}.</p>
      <p class="songName">${this.props.songname}</p>
      <p class=songLength>${this.props.songlength}</p>
      
      <p class="songPrice">$${this.props.price}</p>
      <audio src="/previews/${this.props.uuid}" ref=${this.audio} onEnded=${() => this.setState({playing:false})}></audio>
      ${this.renderDeleteSong()}

      </div>
      `

    }
}
export {Song}
