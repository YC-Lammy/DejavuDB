 """    This program is free software: you can redistribute it and/or modify
    it under the terms of the GNU General Public License as published by
    the Free Software Foundation, either version 3 of the License, or
    (at your option) any later version.

    This program is distributed in the hope that it will be useful,
    but WITHOUT ANY WARRANTY; without even the implied warranty of
    MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
    GNU General Public License for more details.

    You should have received a copy of the GNU General Public License
    along with this program.  If not, see <https://www.gnu.org/licenses/>.
"""

import os, sys


from PySide2.QtWidgets import *
from PySide2.QtCore import *
from PySide2.QtGui import *

from string import ascii_uppercase

########################################################################################################################################################
############################## QAbstractTableModel #####################################################################################################
########################################################################################################################################################

class MyTableModel(QAbstractTableModel): # numpy array model
        # when tableView is rendered, data of each cell would be called through self.data()
        # this method avoid creating widgets for each cell e.g. QtableWiget
        # by storing data in numpy array, fast processing of data could be performed
        # numpy opperation is also supported. 

        def __init__(self, array, headers= None,parent=None):
            super().__init__(parent)

            self.array = array # call current array later through tableWidget.model().array
            self.headers = headers

            self.stack = QUndoStack() # support undo redo function

            # create a dict {'A':'1', 'B':'2', 'C':'3' ...}
            self.di=dict(zip([str((ord(c)%32)-1) for c in ascii_uppercase],ascii_uppercase))

            if '<U' in str(array.dtype) : # '<U' is unicode
                self.numeric = False
            else:
                self.numeric = True

        def formatNumericHeader(self,section):

            section = [i for i in section]
            a =''
            for i in section:
                a += self.di[i]
            return a

        def headerData(self, section: int, orientation: Qt.Orientation, role: int): # fetch the header to GUI
            if role == Qt.DisplayRole:
                if orientation == Qt.Horizontal:
                    if self.headers != None:
                        try:
                            return self.headers[section]  # column headers maybe out of range
                        except :
                            return str(section)
                            #return self.formatNumericHeader(str(section)) # return number instead if out of range
                    else:
                        return str(section)
                        #return self.formatNumericHeader(str(section)) # column
                else:
                    return str(section)  # row

        def columnCount(self, parent=None):
            return len(self.array[0])

        def rowCount(self, parent=None):
            return len(self.array)

        def data(self, index: QModelIndex, role: int): # fetch data to GUI

            if role == Qt.DisplayRole or role == Qt.EditRole:
                row = index.row()
                col = index.column()
                return str(self.array[row][col]) # return value of cell

        def setData(self, index, value, role): # set data would be called everytime user edit cell
            global saved_file # if user modify the array, the array is modified

            if role == Qt.EditRole: # check if func called correctly
                if value:

                    self.stack.push(CellEdit(index, value, self)) # push a new command

                    if value[0] == '=': # '=' indicates to perform a function
                        pass

                    saved_file = False # indicate file modifyed

                    if value.isnumeric() and self.numeric:
                        if 'float' in str(self.array.dtype):
                            value = float(value)
                        elif 'int' in str(self.array.dtype):
                            value = int(value)

                    self.array[index.row()][index.column()] = value # asign new data to array
                    self.update() # update GUI
                    return True
                else:
                    return False # vlue not provided mal function call

        def undo(self):
            self.stack.undo()

        def redo(self):
            self.stack.redo()

        def flags(self, index): # indicate the model's flags
            return Qt.ItemIsEnabled | Qt.ItemIsSelectable | Qt.ItemIsEditable 

    # support for undo redo command
class CellEdit(QUndoCommand): # a new command is pushed to the model stack every time cell is edit

    def __init__(self, index, value, model, *args, **kwargs):
            super().__init__(*args, **kwargs)
            self.index = index # save the cell location
            self.value = value # save the new value 
            self.prev = model.array[index.row()][index.column()] # save the previous value
            self.model = model # a pointer to the model

    def undo(self):
            # set the specific cell to the previous value
            self.model.array[self.index.row()][self.index.column()] = self.prev
            self.model.update()

    def redo(self):
            # set the specific cell to the new value
            self.model.array[self.index.row()][self.index.column()] = self.value
            self.model.update()

