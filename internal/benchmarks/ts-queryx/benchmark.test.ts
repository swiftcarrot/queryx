import { newClient,Model} from "./db";
import * as Benchmark from 'benchmark';

const suite = new Benchmark.Suite();
const c=newClient()
let model: Model


function benchmarkCreate(){
    suite
        .on('complete', (event: { target: any; }) => {
            console.log(String(event.target));
        }).add('create:', async () => {
        model = await c.queryModel().create({
            name: "benchmark",
            title: "title",
            fax: "fax",
            web: "web",
            age: 122,
            righ: true,
            counter: 1222
        }).then();
    }).run()
        .add("insertAll:",async () => {
            await c.queryModel().insertAll([{
                name: "all1benchmark",
                title: "title",
                fax: "fax",
                web: "web",
                age: 122,
                righ: true,
                counter: 1222
            }, {
                name: "all2benchmark",
                title: "title",
                fax: "fax",
                web: "web",
                age: 122,
                righ: true,
                counter: 1222
            }, {
                name: "all3benchmark",
                title: "title",
                fax: "fax",
                web: "web",
                age: 122,
                righ: true,
                counter: 1222
            }])
        }).run()
        .add("find",async () => {
            await c.queryModel().find(model.id)
        }).run()
        .add( "update:", async () => {
            await c.queryModel().where(c.modelID.eq(model.id)).updateAll({
                name: "benchmark-update",
                title: "title",
                fax: "fax",
                web: "web",
                age: 122,
                righ: true,
                counter: 1222
            });
        }).run()
    ;
}
benchmarkCreate()
