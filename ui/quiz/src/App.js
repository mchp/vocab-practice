import React, { Component } from 'react';
import './App.css';

class quiz extends Component {
  constructor(props) {
    super(props);
    this.state = {
      question: "hola", 
      answers: ["hi", "hello"], 
      currentAnswer: "",
      loading: false,
      error: "",
    };
    this.fetchNextUrl = this.props.url; //??
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

  componentDidUpdate() {
    if (this.verdict() && !this.state.loading) {
      console.log("right answer");
      this.fetchNext();
    }
  }

  fetchNext() {
    this.setState({loading: true});
    setTimeout(() => {
      this.setState({
        loading: false,
        currentAnswer: ""});
      document.getElementById("answer").value = ""}, 500)
   /* fetch(this.fetchNextUrl)
      .then(res => res.json())
      .then(
        (result) => {
          this.setState({
            question: result.vocab,
            answers: result.answers,
            currentAnswer: "",
            loading: false,
            error: "",
          });
        },

        (error) => {
          this.setState({
            loading: false,
            error: error,
          })
        });
    */       
  }

  verdict() {
    return this.state.answers.includes(this.state.currentAnswer);
  }

  render() {
    return (
      <div>
        <h1>{this.state.question}</h1>
        <input id="answer" onKeyUp={(e) => this.check(e)}/>
        <button onClick={this.fetchNext}> next </button>
        <div className={this.state.currentAnswer===""?"hidden":""}>
          {this.verdict()?"correct":"incorrect"}
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
