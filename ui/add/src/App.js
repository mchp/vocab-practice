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
      submitted: [],
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
          submitted: [],
        });
      });
  }

  submit() {
    var selected = document.getElementsByClassName("translationCheckbox");
    for (var i=0; selected[i]; i++) {
      if (!selected[i].checked || selected[i].disabled) {
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
          (m) => {
            this.setState({
              loading:false,
              error: m,
              submitted: this.state.submitted.concat(translation),
          })}
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
                var disabled = translation.exists || this.state.submitted.includes(translation.translation);
                return (
                  <div>
                    <input className="translationCheckbox" type="checkbox" name="translation" defaultChecked={disabled} disabled={disabled} value={translation.translation}/> {translation.translation} ({translation.class})
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
