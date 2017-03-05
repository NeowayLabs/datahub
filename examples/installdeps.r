r = getOption("repos")
r["CRAN"] = "http://cran.uk.r-project.org"
options(repos = r)
rm(r)

install.packages("randomForest")
install.packages("rattle")
install.packages("rpart.plot")
install.packages("RColorBrewer")
install.packages("rpart")
