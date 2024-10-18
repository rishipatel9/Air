-- CreateTable
CREATE TABLE "User" (
    "id" SERIAL NOT NULL,
    "Title" TEXT NOT NULL,
    "Desc" TEXT NOT NULL,
    "Completed" BOOLEAN NOT NULL,

    CONSTRAINT "User_pkey" PRIMARY KEY ("id")
);
