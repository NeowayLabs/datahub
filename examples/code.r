library(randomForest)
library(rattle)
library(rpart.plot)
library(RColorBrewer)
library(rpart)

# The train and test data is stored in the ../input directory
train <- read.csv("train.csv")
test  <- read.csv("test.csv")

# We can inspect the train data. The results of this are printed in the log tab below
summary(train)

# Build the decision tree
my_tree_two <- rpart(formula = Survived ~ Pclass + Sex + Age + SibSp + Parch + Fare + Embarked, data=train, method="class")
 
# Visualize the decision tree using plot() and text()
plot(my_tree_two)
text(my_tree_two)
 
# Time to plot your fancified tree
fancyRpartPlot(my_tree_two)
# Make your prediction using the test set
my_prediction <- predict(my_tree_two, test, type = "class")
 
# Create a data frame with two columns: PassengerId &amp; Survived. Survived contains your predictions
my_solution <- data.frame(PassengerId = test$PassengerId, Survived = my_prediction)
 
# Check that your data frame has 418 entries
nrow(my_solution)
 
# Write your solution to a csv file with the name my_solution.csv
write.csv(my_solution, file = "output_solution.csv", row.names = FALSE)
# Here we will plot the passenger survival by class
train$Survived <- factor(train$Survived, levels=c(1,0))
levels(train$Survived) <- c("Survived", "Died")
train$Pclass <- as.factor(train$Pclass)
levels(train$Pclass) <- c("1st Class", "2nd Class", "3rd Class")

png("output_image.png", width=800, height=600)
mosaicplot(train$Pclass ~ train$Survived, main="Passenger Survival by Class",
           color=c("#8dd3c7", "#fb8072"), shade=FALSE,  xlab="", ylab="",
           off=c(0), cex.axis=1.4)
dev.off()
dim(train)
