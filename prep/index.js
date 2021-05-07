const fs = require("fs");
const path = require("path");

const txt = fs.readFileSync(path.join(__dirname,"doc_dump.txt")).toString();
const r = txt.split("\n").map(x=>x.trim());

fs.mkdirSync(path.join(__dirname,"out"));

r.forEach(line => {
    const parts = line.split("\t").map(x=>x.trim());
    const id = parts[0];
    const abstract = parts[3];
    fs.writeFileSync(path.join(__dirname,"out",id+".txt"),abstract);
});
