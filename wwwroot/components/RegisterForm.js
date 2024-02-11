import {Component,h,createRef} from '/js/preact.js';
import htm from '/js/htm.js'

const html = htm.bind(h)

class RegisterForm extends Component {
    displayname = createRef("displaynametxt")
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
        let displayname = this.displayname.current.value
        let username = this.username.current.value
        let password = this.password.current.value
        if(username == "") {
            this.setState({error:"Displayname cannot be empty"});
            return
        }
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
        "DisplayName": displayname,
        "Username": username,
        "PasswordHash": password
    });
    
    var requestOptions = {
      method: 'POST',
      headers: myHeaders,
      body: raw,
      redirect: 'follow'
    };

    let that = this
    
    fetch("/api/register", requestOptions)
      .then(response => {
        if (response.status === 409) {     
            throw new Error("Account already exists")
        }
        response.text()
        })
        .then(result => { 
        window.location.pathname = "/login.html" 
        })
        .catch(error => {
            this.setState({error: error.message})
        });
    }
    render() {
      return html `
      <div class="formDiv">
      <div class="formBox">
      ${this.rendererrorbanner()}
      <h1 class="formTitle">Register</h1>
      
      <div class="formelementdiv">
      <input id="displaynametxt" ref=${this.displayname}  placeholder="Display Name"/>
      </div>

      <div class="formelementdiv">
      <input id="usernametxt" ref=${this.username}  placeholder="Username"/>
      </div>
      
      <div class="formelementdiv">
      <input id="passwordtxt" ref=${this.password} type="password" placeholder="Password"/>
      </div>

      <div class="formelementdiv">
      <button id=submit onClick=${(e) => this.buttonclick(e)}>Register</button>
      </div>

      </div>
      </div>
      `

    }
}
export {RegisterForm};