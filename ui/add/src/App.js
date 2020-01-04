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
        <div className="button">
          <a onClick={this.lookup}>Look up </a>
        </div>
        <div className={this.state.translations.length>0?"":"hidden"}>
          <div id="results">
            {
              this.state.translations.map(translation => {
                var disabled = translation.exists || this.state.submitted.includes(translation.translation);
                if (disabled) {
                  return (
                    <div>
                      <input className="translationCheckbox" type="checkbox" name="translation" checked disabled value={translation.translation}/> {translation.translation} ({translation.class})
                    </div>
                )}
                return (
                  <div>
                    <input className="translationCheckbox" type="checkbox" name="translation" value={translation.translation}/> {translation.translation} ({translation.class})
                  </div>
                )
              })  
            }
          </div>
        <div className="button">
          <a onClick={this.submit}>Submit!</a>
        </div>
        </div>
        <div className={this.state.loading?"":"hidden"}>
          loading...
        </div>
        <div className="error">
          {this.state.error}
        </div>
      </div>
    );
  }
}

export default add;
