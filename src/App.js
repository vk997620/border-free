import React, { useEffect , useState , PureComponent } from "react";
import { LineChart, Line, Tooltip , CartesianGrid, XAxis, YAxis } from 'recharts';

import './App.css';
class App extends React.Component {
  constructor(props) {
    super(props);
    this.state = {
      error: null,
      isLoaded: false,
      items: []
    };
    
  }

  com  = function compare(a, b){

    if(a.id < b.id) return -1;
    else if(a.id == b.id) return 0;
    else return 1;
  }



  data = [];
  
  componentDidMount() {
    fetch("https://0pz3nv144j.execute-api.ap-south-1.amazonaws.com/test/")
      .then(res => res.json())
      .then(
        (result) => {
          this.setState({
            isLoaded: true,
            items: result
          });
        },
        (err)=>{
          this.setState({
            isLoaded:true,
            error:err
          })
        }
      );
  }

  render() {
    const { error, isLoaded, items } = this.state;
    if (error) {
      return <div>Error: {error.message}</div>;
    } else if (!isLoaded) {
      return <div>Loading...</div>;
    } else {
      // var body = JSON.parse(items.body);

      var body = JSON.parse(items.body);

      for(let i = 0 ; i < body.hits.hits.length ; i++){
          let item = body.hits.hits[i]._source;
        //  console.log(i , item.Cupcakes , item.Year , item.Month)
        this.data.push({id:body.hits.hits[i]._id , Year : item.Year , Month:item.Month , Date:item.Year+"-"+item.Month, Cupcakes : item.Cupcakes})
      }

      this.data.sort(this.com)

      for(let i = 0 ; i < this.data.length ; i++){
        console.log(this.data[i].id , this.data[i].Year , this.data[i].Month , this.data[i].Cupcakes)
      }



      return (
        <div>
        <LineChart width={1350} height={450} margin={{ top: 150, right: 30, left: 20, bottom: 5 }} data={this.data}>
        <Line type="monotone" dataKey="Cupcakes" stroke="#8884d8" dot={false}/>
        <XAxis dataKey="Date" tickCount={5}/>
        <CartesianGrid strokeDasharray="3 3" />
        <YAxis/>
        <Tooltip />
        </LineChart>
        </div>     

      )
    }
  }
}  

export default App;

