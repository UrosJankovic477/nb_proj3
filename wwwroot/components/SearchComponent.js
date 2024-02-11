import {Component,h, createRef} from '/js/preact.js';
import htm from '/js/htm.js'
import { AlbumDisplay } from '/components/AlbumDisplayComponent.js';


const html = htm.bind(h)

class Search extends Component {

    formref = createRef("formref")

    async componentDidMount() {
        

        let maxprice = 0
        let genres = []

        
        
        var requestOptions = {
          method: 'GET',
          redirect: 'follow'
        };

        try
        {
            let response = await fetch("/api/getSearchOptions", requestOptions)
            let result = await response.json()
            maxprice = result.MaxPrice
            genres = result.Genres
        }
        catch (error) {
            console.log(error)
        }


        let assortment = {genres:genres,minprice:0,maxprice:maxprice}
        this.setState({globalAssortmentInfo:assortment, choice: {
            genres: [],
            maxprice: maxprice,
            query_string: "",
        }})
     }

    renderGenres() {
        let genrelist = []
        let info = this.state.globalAssortmentInfo
        if(!this.state.globalAssortmentInfo) {return null}
        info.genres.forEach((genre => {
            genrelist.push(html `<div class=checkboxdiv><input id="check_${genre}" type=checkbox name="genres" value="${genre}"/><label for="check_${genre}">${genre}</label></div>`)
        }))
        return genrelist
    }
    renderMinMaxPrice() {
        let info = this.state.globalAssortmentInfo
        if(!this.state.globalAssortmentInfo) {return null}
        const minprice=info.minprice.toFixed(2)
        const maxprice=info.maxprice.toFixed(2)
        return html`<div class="minmaxdiv"><p>$${minprice}</p><input id="slider" type="range" step="0.25" min="${minprice}" max=${maxprice}/><p>$${maxprice}</p></div>`
    }

    changeOptions(e) {
        let form = this.formref.current
        let checkedGenres = form.querySelectorAll("input:checked")
        let genres = []
        checkedGenres.forEach(x => {
            genres.push(x.value)
        })
        let query_string = form.querySelector("input#queryString").value
        let slider = form.querySelector("input#slider").value

        let assortment = JSON.parse(JSON.stringify(this.state.globalAssortmentInfo))
        assortment.choice = {
            genres: genres,
            maxprice: Number(slider),
            query_string: query_string,
        }
        this.setState(assortment)
        
        
    }

    renderAlbumDisplay() {
        if(this.state.choice === undefined) {
            return null
        }
        return html `<${AlbumDisplay} id="albumDisplay" clickaction=${this.props.clickaction}
        genres="${this.state.choice.genres.join("_")}" 
        maxprice="${this.state.choice.maxprice}" 
        query_string="${this.state.choice.query_string}"/>`
    }
    
    render() {
      return html `
      <div class="storeView">
      <div class="smallerDiv">
      <form id="searchForm" ref=${this.formref} onChange=${(e) => {this.changeOptions(e)}}>
      <input id="queryString" type="text" onInput=${(e) => {this.changeOptions(e)}} placeholder="Search by artist or song name"/>
      ${this.renderMinMaxPrice()}
      ${this.renderGenres()}
      </form>
      </div>
      <div class="biggerDiv">
      ${this.renderAlbumDisplay()}
      </div>
      </div>
      `
    }

}
export {Search}
