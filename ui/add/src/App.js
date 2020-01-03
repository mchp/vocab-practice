import React, { Component } from 'react';
import './App.css';

class add extends Component {
  constructor(props) {
    super(props);
    this.state = { 
      vocab: "",
      translations: [],
      selected: "",
      loading: false,
      error: "",
    }
    this.lookup = this.lookup.bind(this);
    this.submit = this.submit.bind(this);
  }
  
  lookup() {
    this.setState({loading: "true", error: ""});
    var vocab = document.getElementById("vocab").value;
    fetch("/lookup?vocab=" + vocab)
      .then(res => res.json())
      .then(result => {
        this.setState({
          loading: false,
          vocab: vocab,
          translations: result ? result : [],
          selected: "",
          error: result ? "" : "No result found",
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
      })
        .then(res => res.text())
        .then(
          (m) => {this.setState({loading:false, error : m})}
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
        <div className={this.state.loading?"":"hidden"}>
          loading...
        </div>
        {this.state.error}
      </div>
    );
  }
}

export default add;
