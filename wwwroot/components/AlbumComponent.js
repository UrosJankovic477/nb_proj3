import {Component,h} from '/js/preact.js';
import htm from '/js/htm.js'


const html = htm.bind(h)

class Album extends Component {
    renderArtist() {
        if(!this.props.artistname) {
            return null
        }
        return html`<a href="/artist/${this.props.artistname}" class="artistname">${this.props.artistdisplayname}</p>`
    }

    

    render() {
      return html `
      <div class="albumDiv hovershadow albumInfo" onMouseUp=${(e) => this.props.clickaction(e,this.props.uuid)}>
      <img src="/albumart/${this.props.uuid}"></img>
      <a href="/album/${this.props.uuid}" class="albumname">${this.props.albumname}</p>
      ${this.renderArtist()}
      <p class="price">$${this.props.price}</p>
      </div>
      `

    }
}
export {Album}
