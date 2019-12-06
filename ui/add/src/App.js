import React, { Component } from 'react';
import './App.css';

class add extends Component {
  constructor(props) {
    super(props);
    this.state = { 
      vocab: "bonita",
      translations: [{translation:"pretty", class:"adjective"}, {translation:"beautiful", class:"adjective"}],
      selected: "",
      loading: false,
      error: false,
    }
    this.lookup = this.lookup.bind(this);
    this.submit = this.submit.bind(this);
  }
  
  lookup() {
    this.setState({loading: "true"});
    var vocab = document.getElementById("vocab").value;
    fetch("/lookup?vocab=" + vocab)
      .then(res => res.json())
      .then(result => {
        this.setState({
          loading: false,
          vocab: vocab,
          translations: result,
          selected: "",
        });
      });
  }

  submit() {
    var selected = document.getElementsByClassName("translationCheckbox");
    for (var i=0; selected[i]; i++) {
      if (!selected[i].checked) {
        continue;
      }
      var translation = selected[i].value;
      fetch("/input", {
        method: 'POST',
        headers: {'Content-Type': 'application/x-www-form-urlencoded'},
        body: 'vocab=' + this.state.vocab + '&translation=' + translation,
      }).then(
        () => {this.setState({loading: false})},
        (error) => {this.setState({loading:false, error : error.message})}
      )
  }};
  
  render() {
    return (
      <div>
        <input id="vocab" />
        <button onClick={this.lookup}>lookup</button>
        <div id="results">
          <div>
            {
              this.state.translations.map(translation => {
                return (
                  <div>
                    <input className="translationCheckbox" type="checkbox" name="translation" value={translation.translation}/> {translation.translation} ({translation.class})
                  </div>
              )})
            }
          </div>
          <button onClick={this.submit}>submit</button>
        </div>
      </div>
    );
  }
}

export default add;
