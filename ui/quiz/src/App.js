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
    var currentAnswer = document.getElementById("answer").value;
    if (currentAnswer === this.state.currentAnswer) {
      return;
    }
    this.setState({
      currentAnswer: currentAnswer,
    });
  }

  componentDidMount() {
    this.fetchNext();
  }

  componentDidUpdate(prevProps, prevState) {
    if (prevState && this.state.currentAnswer === prevState.currentAnswer) {
      return;
    }
    if (this.verdict() && !this.state.loading) {
      console.log("right answer");
      this.postPass();
    }
  }

  postPass() {
    this.setState({loading: true});
    fetch('/pass', {
      method: 'POST',
      headers: {'Content-Type': 'application/x-www-form-urlencoded'},
      body: 'vocab=' + this.state.question + 
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
      .then(res => res.json())
      .then(
        (result) => {
          this.setState({
            question: result.vocab,
            answers: result.translations,
            currentAnswer: "",
            loading: false,
            error: "",
          });
          document.getElementById("answer").value = "";
        },

        (error) => {
          this.setState({
            loading: false,
            error: error.message,
          })
        }); 
  }

  verdict() {
    return this.state.currentAnswer !== "" && this.state.answers.some(a => a.translation === this.state.currentAnswer);
  }

  renderAnswer(answer, time) {
    return <div> {answer} (last tested {time}) </div>
  }

  render() {
    return (
      <div>
        <h1>{this.state.question}</h1>
        <input id="answer" onKeyUp={(e) => this.check(e)}/>
        <button onClick={this.fetchNext}> next </button>
        <div className={this.state.currentAnswer===""?"hidden":""}>
          {this.verdict()?"correct":"incorrect"}
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
