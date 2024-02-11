import {Component,h,createRef} from '/js/preact.js';
import htm from '/js/htm.js'


const html = htm.bind(h)

class LoginForm extends Component {
    username = createRef("usernametxt")
    password = createRef("passwordtxt")
    constructor(props) {
    
       
      super(props);
      this.state = {error:props.error}
      
      
    }

    rendererrorbanner() {
      
      if(!this.state.error)  {
        return null
      }
      return html`<div class=errordiv>${this.state.error}</div>`

    }

    buttonclick(e) {
      let username = this.username.current.value
      let password = this.password.current.value
      if(username == "") {
        this.setState({error:"Username cannot be empty"});
        return
      }
      if(password == "") {
        this.setState({error:"Password cannot be empty"});
        return
      }

      var myHeaders = new Headers();
      myHeaders.append("Content-Type", "application/json");
      
      var raw = JSON.stringify({
        "Username": username,
        "PasswordHash": password
      });

      var requestOptions = {
        method: 'POST',
        headers: myHeaders,
        body: raw,
        redirect: 'follow'
      };

      fetch("/api/login", requestOptions)
        .then(response => response.text())
        .then(result => {
          document.cookie = `token=${JSON.parse(result)}`
          window.location.pathname = "/index.html"
        })
        .catch(error => this.setState({error: error}))
    }
    render() {
      return html `
      <div class="formDiv">
      <div class="formBox">
      ${this.rendererrorbanner()}
      <h1 class="formTitle">Login</h1>
      <div class="formelementdiv">
      <input id="usernametxt" ref=${this.username}  placeholder="Username"/>
      </div>
      
      <div class="formelementdiv">
      <input id="passwordtxt" ref=${this.password} type="password" placeholder="Password"/>
      </div>

      <div class="formelementdiv">
      <button id=submit onClick=${(e) => this.buttonclick(e)}>Login</button>
      </div>

      </div>
      </div>
      `

    }
}
export {LoginForm};
