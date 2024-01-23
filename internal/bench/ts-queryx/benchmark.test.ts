import { newClient,Model} from "./db";
import * as Benchmark from 'benchmark';

const suite = new Benchmark.Suite();
const c=newClient()
let model: Model


function benchmarkCreate(){
    suite
        .on('complete', (event: { target: any; }) => {
            console.log(String(event.target));
        }).add('create:', () => {
        model =  c.queryModel().create({name: "benchmark", title: "title", fax: "fax", web: "web", age: 122, righ: true, counter: 1222}).then();
        }).run()
        .add("insertAll:",()=>{
         let models =  c.queryModel().insertAll([{name:"all1benchmark",title:"title",fax:"fax",web:"web",age:122,righ: true,counter:1222},{name:"all2benchmark",title:"title",fax:"fax",web:"web",age:122,righ: true,counter:1222},{name:"all3benchmark",title:"title",fax:"fax",web:"web",age:122,righ: true,counter:1222}])
        }).run()
        .add("find",()=>{
         let m =   c.queryModel().find(model.id)
        }).run()
        .add( "update:", ()=>{
         let u =  c.queryModel().where(c.modelName.eq(model.name)).updateAll({name:"benchmark-update",title:"title",fax:"fax",web:"web",age:122,righ: true,counter:1222});
        }).run()
    ;
}
benchmarkCreate()


