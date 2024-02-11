import {Component,h} from '/js/preact.js';
import htm from '/js/htm.js'


const html = htm.bind(h)

class Header extends Component {

    async getLoggedIn() {
        if (document.cookie) {
            var myHeaders = new Headers();
            myHeaders.append("Content-Type", "application/json");
            let token = document.cookie.replace("token=", "");

            var raw = JSON.stringify(token);

            var requestOptions = {
              method: 'POST',
              headers: myHeaders,
              body: raw,
              redirect: 'follow'
            };
            let ret_val = undefined
            try {
                let response = await fetch("/api/getLoggedInUsersData", requestOptions)
                let result = await response.json()
                let user = result
                return {
                    cookie:document.cookie,loggedin:true, userInfo: {
                        displayname:user.DisplayName,
                        username:user.Username, 
                        balance:user.Balance, 
                        itemsincart: user.CartCount
                }
            }
                
            } catch (Fehler) {
                return Fehler
            }
            
            
            
        }
        return {cookie:null,loggedin:false}
    }
        
            componentDidMount() {
                this.getLoggedIn().then(login => {
                    this.setState(login)
                    if(this.props.requirelogin && !login.loggedin) {
                        window.location.hash = "loginrequired"
                        window.location.pathname = "/login.html"
                    }
                })  
                

            }

    rendernavBarButtons() {
        if(!this.state.loggedin) {
            return html`<a class=undecorated href=/login.html>Login</a><a class=undecorated href=/register.html>Register</a>`
        }
        return html`<a class=undecorated href=/myaccount.html>${this.state.userInfo.displayname} ($${this.state.userInfo.balance.toFixed(2)})</a><a class=undecorated href=/cart.html>Cart(${this.state.userInfo.itemsincart})</a>`

    }

    render() {

        return html `
        <div class="navbar">
        ${this.rendernavBarButtons()}
        </div>
        `

    }
}
export {Header}
