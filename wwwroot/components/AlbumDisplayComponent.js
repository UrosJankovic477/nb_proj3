import {Component,h} from '/js/preact.js';
import {Album} from '/components/AlbumComponent.js'
import htm from '/js/htm.js'

let count = 20
let page = 0

const html = htm.bind(h)

class AlbumDisplay extends Component {

    

    async addAlbum(e) {
      let name = prompt("Enter album name")
      var myHeaders = new Headers();
            myHeaders.append("Content-Type", "application/json");
            let token = document.cookie.replace("token=", "");

            var raw = JSON.stringify({
              Token: token,
              Title: name,
            });
            

            var requestOptions = {
              method: 'POST',
              headers: myHeaders,
              body: raw,
              redirect: 'follow'
            };
            let ret_val = undefined
            try {
                let response = await fetch("/api/addAlbum", requestOptions)
                let result = await response.json()
                window.location.pathname = `/album/${result}`
                
            }  
            catch (Fehler) {
                return Fehler
            }
    }

    async componentDidMount() {
      let token = document.cookie.replace("token=", "")
      if(window.location.pathname.includes("/artist/")){
        var myHeaders = new Headers();
        myHeaders.append("Content-Type", "application/json");
        let uuid = window.location.pathname.replace("/artist/", "")
       
        var raw = JSON.stringify(
          uuid
        );
       
        var requestOptions = {
          method: 'POST',
          headers: myHeaders,
          body: raw,
          redirect: 'follow'
        };

 
        let response = await fetch("/api/getArtistsAlbums", requestOptions)
        let result = await response.json()
        this.setState({
          albumlist: result,
        })

      }

      else if(window.location.pathname === "/cart.html")
      {
        
        var myHeaders = new Headers();
        myHeaders.append("Content-Type", "application/json");
       
        var raw = JSON.stringify(
          token
        );
       
        var requestOptions = {
          method: 'POST',
          headers: myHeaders,
          body: raw,
          redirect: 'follow'
        };

 
        let response = await fetch("/api/getCart", requestOptions)
        let result = await response.json()
        this.setState({
          albumlist: result,
        })
      }
    }

    async componentDidUpdate() {
      if (window.location.pathname !== "/" ) {
        return
      }
      let props_stringy = JSON.stringify(this.props)
      let old_props_stringy = JSON.stringify(this.state.props)
      if (old_props_stringy === props_stringy) {
        return
      }
      

      
      
        var myHeaders = new Headers();
        myHeaders.append("Content-Type", "application/json");
        let maxprice = 0
        let genres = []
        let query_string = ""
        if (this.props.maxprice !== undefined) {
          maxprice = Number(this.props.maxprice)
        }
        if (this.props.genres) {
          genres = this.props.genres.split("_")
        }
        if (this.props.query_string) {
          query_string = this.props.query_string
        }
        var raw = JSON.stringify({
          "Count": count,
          "Page": page,
          "MaxPrice": maxprice,
          "Genres": genres,
          "Query": query_string,
        });
       
        var requestOptions = {
          method: 'POST',
          headers: myHeaders,
          body: raw,
          redirect: 'follow'
        };

 
        let response = await fetch("/api/search", requestOptions)
        let result = await response.json()
        this.setState({
          albumlist: result,
          props: this.props,
        })

    }

    renderAlbums() {
        let renderedlist = []
        if(!this.state.albumlist) {return null}
        this.state.albumlist.forEach((album) => {
          renderedlist.push(html`<${Album} clickaction=${this.props.clickaction} albumname="${album.Title}" artistname="${album.ArtistUsername}" artistdisplayname="${album.ArtistDisplayname}" price="${album.Price.toFixed(2)}" uuid="${album.UUID}"/>`)
       })
       return renderedlist
    }

    renderAdditionalButton() {
      if(this.props.source == "creations") {
        return html `<button class="addButton albumDiv" onClick=${(e) => this.addAlbum(e)}>+</button>`
      }

     
    }

    render() {
      return html `
      <div class="albumDisplayBox">
      ${this.renderAlbums()}
      ${this.renderAdditionalButton()}
      </div>
      `

    }
}
export {AlbumDisplay}
