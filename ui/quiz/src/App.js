import React, { Component } from 'react';
import './App.css';

class quiz extends Component {
  constructor(props) {
    super(props);
    this.state = {
      question: "", 
      answers: [], 
      currentAnswer: "",
      loading: false,
      error: "",
    };
    this.check = this.check.bind(this);
    this.verdict = this.verdict.bind(this);
    this.fetchNext = this.fetchNext.bind(this);
  }

  check(e) {
    if (e.keyCode !== 13) {
      return;
    }
    var currentAnswer = document.getElementById("answer").value.toLowerCase();
    if (currentAnswer === this.state.currentAnswer) {
      this.fetchNext();
    } else {
      this.setState({
        currentAnswer: currentAnswer,
    })};
  }

  componentDidMount() {
    this.fetchNext();
  }

  componentDidUpdate(prevProps, prevState) {
    if (prevState && this.state.currentAnswer === prevState.currentAnswer) {
      return;
    }
    if (this.verdict() && !this.state.loading) {
      this.postPass();
    }
  }

  postPass() {
    this.setState({loading: true});
    fetch('/pass', {
      method: 'POST',
      headers: {'Content-Type': 'application/x-www-form-urlencoded'},
      body: 'vocab=' + this.state.question.toLowerCase() + 
            '&translation=' + this.state.currentAnswer
    }).then(() => {
        this.setState({
          loading: false,
        });
    },
      (error) => {
        this.setState({
          loading: false,
          error: error.message,
        })
      });
  }

  fetchNext() {
    this.setState({loading: true});
    fetch('/next')
      .then(res => {
        if (res.status == 200) {
          res.json().then(
            (result) => {
              this.setState({
                question: result.vocab,
                answers: result.translations,
                currentAnswer: "",
                loading: false,
                error: "",
              });
            document.getElementById("answer").value = "";
          },)
        } else {
          res.text().then( 
            (error) => {
              this.setState({
                loading: false,
                error: error
              })
        })
      }
    });
  }

  verdict() {
    return this.state.currentAnswer !== "" && this.state.answers.some(a => a.translation.toLowerCase() === this.state.currentAnswer.toLowerCase());
  }

  renderAnswer(answer, time) {
    var testedTimeString = "never tested";
    if (time != "0001-01-01T00:00:00Z") {
      testedTimeString = "last tested " + time;
    }
    return (
    <div>
      <div className={this.state.currentAnswer===answer?"checked":""}> 
        {answer}
      </div>
      <div className="testTime">{testedTimeString}</div>
    </div>)
  }

  render() {
    var classNames = require('classnames');
    var correctClass = classNames("verdict", "correct", {'hidden': !this.verdict()});
    var incorrectClass = classNames("verdict", "incorrect", {'hidden': this.verdict()});
    return (
      <div>
        <div id="card">
          <h1>{this.state.question}</h1>
          <input id="answer" autocapitalize="none" autocomplete="off" onClick={() => {document.getElementById("vocab").select()}} onKeyUp={(e) => this.check(e)}/>
        </div>
        <div id="skip" onClick={this.fetchNext}> skip >> </div>
        <div className={this.state.currentAnswer===""?"hidden":""}>
          <span className={correctClass}>Correct!</span>
          <span className={incorrectClass}>Incorrect :(</span>
          {this.state.answers.map(a => {
            return this.renderAnswer(a.translation, a.lastTested);
          })}
        </div>
        
        <div className={this.state.loading?"":"hidden"}>
          loading...
        </div>
        {this.state.error}
      </div>
    );
  }
}

export default quiz;
